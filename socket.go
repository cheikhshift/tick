package tick

import (
	"log"
	"net/http"
	"time"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
	"github.com/cheikhshift/db"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 15 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 5) / 10

	// Poll file for changes with this period.
	filePeriod = 1 * time.Second
)

var (
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  128,
		WriteBufferSize: 128,
	}
	Key string
	DB db.DB
)



func SetDb(db  db.DB) {
	DB = db
}

// encrypt string to base64 crypto using AES
func Encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return "Error"
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return "Error"
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func sResponse(v interface{}) string {
					data,_ := json.Marshal(&v)
					return string(data)
}

// decrypt from base64 to decrypted string
func Decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}


func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		} 
	}
}

func writer(ws *websocket.Conn,collection string, id string ) {
	//lastError := ""
	//pingTicker := time.NewTicker(pingPeriod)
	/*fileTicker := time.NewTicker(filePeriod)
	defer func() {
		//pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()*/
	//var temp string
	//	temp = ""

		
				//ws.SetWriteDeadline(time.Now().Add(writeWait))
				//get obj
				var obj  interface{}
				
				DB.MoDb.C(collection).Find(bson.M{"_id" : bson.ObjectIdHex(id)}).One(&obj)
				//newResponse := sResponse(obj)
				
				if err := ws.WriteMessage(websocket.TextMessage, []byte(sResponse(obj))); err != nil {
					fmt.Println("Socket Error : " ,err.Error())
				}
				time.Sleep(filePeriod)
				writer(ws, collection,id)
	
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}


	path := strings.Split(Decrypt([]byte(Key), r.FormValue("token")), ",")

	go writer(ws, path[0], path[1])
	//go reader(ws)
}





