package chat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handleResponse(w http.ResponseWriter, a interface{}) {
	// fmt.Println("Handle `handleResponse`")
	resp := response{
		Response: a,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println(string(respJSON))
	fmt.Fprintln(w, string(respJSON))
}

func handleError(w http.ResponseWriter, err error, code int) {
	// fmt.Println("Handle `handleError`")
	resp := errorResp{
		Error: err.Error(),
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Println(string(respJSON))
	http.Error(w, string(respJSON), code)
}

func (chatExpl *ChatExplorer) HandlerUserAdd(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "dfg")
	// http.Error(w, "sdf", 201)
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	user := User{}
	err := json.Unmarshal(body, &user)

	if err != nil {
		handleError(w, err, http.StatusConflict)
	}

	id, err := chatExpl.addUser(user.Username)
	if err != nil {
		handleError(w, err, http.StatusConflict)
		return
	}

	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}
	handleResponse(w, resp)
}

func (chatExpl *ChatExplorer) HandlerChatAdd(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	chat := Chat{}
	err := json.Unmarshal(body, &chat)

	if err != nil {
		handleError(w, err, http.StatusConflict)
		return
	}

	id, err := chatExpl.addChat(chat.Name, chat.Users)
	if err != nil {
		handleError(w, err, http.StatusConflict)
		return
	}

	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}
	handleResponse(w, resp)
}

func (chatExpl *ChatExplorer) HandlerMessageAdd(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	message := Message{}
	err := json.Unmarshal(body, &message)

	if err != nil {
		handleError(w, err, http.StatusConflict)
		return
	}

	id, err := chatExpl.addMessage(message.Chat, message.Author, message.Text)
	if err != nil {
		handleError(w, err, http.StatusConflict)
		return
	}

	resp := struct {
		ID int64 `json:"id"`
	}{
		ID: id,
	}
	handleResponse(w, resp)
}

func (chatExpl *ChatExplorer) HandlerChatsGet(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	user := User{}
	err := json.Unmarshal(body, &user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}
	fmt.Println(user)

	data, err := chatExpl.getChats(user.ID)
	if err != nil {
		handleError(w, err, http.StatusConflict)
	}

	resp := struct {
		Chats []interface{} `json:"chats"`
	}{
		Chats: data,
	}
	handleResponse(w, resp)
}

func (chatExpl *ChatExplorer) HandlerMessagesGet(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	chat := Chat{}
	err := json.Unmarshal(body, &chat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
	}
	fmt.Println(chat)

	data, err := chatExpl.getMessages(chat.ID)
	if err != nil {
		handleError(w, err, http.StatusConflict)
	}

	resp := struct {
		Messages []interface{} `json:"messages"`
	}{
		Messages: data,
	}
	handleResponse(w, resp)
}
