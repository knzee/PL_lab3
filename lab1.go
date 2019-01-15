package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strconv"
import "math/rand"	
import "time"
import "strings"
import "unicode"

func client(ip string) {
	conn, _ := net.Dial("tcp", ip)

	s_hash := get_hash()
	
	for i:=0; i<3; i++ {

		s_key := get_key()
			
		fmt.Print("Sending: hash- ",s_hash," key- ",s_key +"\n")
		
		fmt.Fprintf(conn, s_hash + " " + s_key + "\n")

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Receiving: "+ message)
		if (message[:len(message)-1] == next_session_key(s_hash,s_key)) {
			fmt.Print("<<<Key match>>>" + "\n\n")
		} else {
			break
		}
	}
}

func server(port string) {

  fmt.Println("Launching server at port: " + port + "\n")

  ln, _ := net.Listen("tcp", ":" + port)

  for {
	conn, _ := ln.Accept()
	go handleRequest(conn)
  }
}

func handleRequest(conn net.Conn) {

	for i:=0; i<3;i++ {
		message, _ := bufio.NewReader(conn).ReadString('\n')

		fmt.Print("Incoming data: \n")
		
		a:= strings.Split(message," ")
		a[1]=(a[1])[:len(a[1])-1]
		fmt.Print("hash: ",a[0]," key: ",a[1],"\n\n")
		
		newmessage := next_session_key(a[0],a[1])

		conn.Write([]byte(newmessage + "\n"))
	}
	
	conn.Close()
	fmt.Print("Connection closed" + "\n\n")

}

func get_key() string{
	result := ""
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=1; i < 11; i++ {
		result += string(strconv.Itoa(int(9*r.Float64() + 1))[0])
	}
	return result
}

func get_hash() string{
	li := ""
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0; i < 5; i++ {
		li += strconv.Itoa(int(6*r.Float64() + 1))
	}
	return li
}

func calc_hash(session_key string,val int) string{
	result := ""
	if (val == 1) {
		t1,_ := strconv.Atoi(session_key[0:5])
		temp := "00" + strconv.Itoa(t1 % 97)
		return temp[len(temp)-2:]
	} else if (val == 2) {
		for i:= 1; i < len(session_key); i++ {
			result+= string(session_key[len(session_key)-i])
		}
		return result + session_key[0:5]
	} else if (val == 3) {
		return session_key[(len(session_key)-5):] + session_key[0:5]
	} else if (val == 4) {
		num := 0
		for i:=1; i<9; i++ {
			temp,_ := strconv.Atoi(string(session_key[i]))
			num += temp + 41
		}
		return strconv.Itoa(num)
	} else if (val == 5) {
		num := 0
		for i:=0; i<len(session_key); i++ {
			ch := string(([]rune(session_key)[i])^43 )
			if !unicode.IsDigit([]rune(ch)[0]) {
				ch = strconv.Itoa(int([]rune(ch)[0]))
			}
			temp,_ := strconv.Atoi(ch)
			num += temp
		}
		return strconv.Itoa(num)
	} else {
		temp,_:= strconv.Atoi(session_key)
		return strconv.Itoa(temp + val)
	}
}


func next_session_key(hash string,session_key string) string{
	result := 0
	for i:=0; i<len(hash); i++ {
		temp,_ := strconv.Atoi(calc_hash(session_key, int([]rune(hash)[i])))
		result += temp
	}
	r:= strings.Repeat("0",10) + strconv.Itoa(result)[0:10]
	return r[len(r)-10:]
}


func main() {
	
    args := os.Args[1:]

	if (len(args[0]) > 4) {
		test, _ := strconv.Atoi(args[2])
		if ((args[1]=="-n")&&(test > 0)) {
			for i:= 0; i < test; i++ {
				client(args[0])
			}
		} else {
			fmt.Println("Не указано количество клиентов")
		}

	} else {
		server(args[0])
	}

	var kostil string
	fmt.Fscan(os.Stdin, &kostil) 
}
