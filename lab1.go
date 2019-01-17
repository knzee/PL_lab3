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
	
	s_key := get_key()
	
	var answer [15]byte
	
	conn.Write([]byte(s_hash))
	
	bufio.NewReader(conn).Read(answer[0:15])
	fmt.Print(string(answer[0:15]))
	
	for i:=0; i<3; i++ {
		var new_key [10]byte

		conn.Write([]byte(s_key))
		
		fmt.Print("Sending: hash- ",s_hash," key- ",s_key +"\n")
		
		bufio.NewReader(conn).Read(new_key[0:10])
		message := string(new_key[0:10])
		
		fmt.Print("Receiving: "+ message)
		
		if (message == next_session_key(s_hash,s_key)) {
			fmt.Print("\n<<<Key match>>>" + "\n\n")
			s_key = next_session_key(s_hash,s_key)
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
	var byte_hash [5]byte
	bufio.NewReader(conn).Read(byte_hash[0:5])
	
	hash := string(byte_hash[0:5])
	
	conn.Write([]byte("Hash recieved!\n"))
	
	for i:=0; i<3;i++ {
		var byte_key [10]byte
		bufio.NewReader(conn).Read(byte_key[0:10])
		message:= string(byte_key[0:10])

		fmt.Print("Incoming data: \n")
		
		fmt.Print("hash: ",hash," key: ",message,"\n\n")
		
		newmessage := next_session_key(hash,message)

		conn.Write([]byte(newmessage))
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
		return result + string(session_key[0])
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
		index,_ := strconv.Atoi(string(hash[i]))
		temp,_ := strconv.Atoi(calc_hash(session_key, index))
		result += temp
	}
	l:= len(strconv.Itoa(result))
	if (len(strconv.Itoa(result)) > 10 ) {l = 10}
	r:= strings.Repeat("0",10) + strconv.Itoa(result)[0:l]
	return r[len(r)-10:]
}


func main() {

    args := os.Args[1:]

	if (len(args[0]) > 4) {
		test, _ := strconv.Atoi(args[2])
		if ((args[1]=="-n")&&(test > 0)) {
			for i:= 0; i < test; i++ {
				go client(args[0])
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
