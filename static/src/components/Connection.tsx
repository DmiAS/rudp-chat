import React, {useEffect} from "react"
import { Msg, User } from "../interfaces/core"
import {Loader} from "./Loader"
import {UserComp} from "./UserComp"
import axios from "axios"
import { Chat } from "./Chat"


type ConnectionProps = {
    users: User[]
    isLoading: boolean
}

//! МОКИ. Тут должно быть получение сообщений из сокета
// const msgs: Msg[] = [{text: "Привет"}, {text: "1234123412341234123412341234", fromMe: true}]

export const Connection: React.FC<ConnectionProps> = ({ users, isLoading }) => {
    const [connected, setConnected] = React.useState(false)
    const [isChatLoading, setIsChatLoading] = React.useState(false)
    // const [chosen, setChosen] = React.useState({name: '', id: -1})
    const [address, setAddress] = React.useState("")
    const [messages, setMessages] = React.useState<Msg[]>([])
    const [msgSocket, setSocket] = React.useState<WebSocket>()
    const [fileSocket, setFileSocket] = React.useState<WebSocket>()
    const [name, setName] = React.useState<string>("")

    useEffect( () => {
        let new_uri = "ws://" + window.location.host + "/ws/chat/thread"
        console.log(new_uri)
        let ws = new WebSocket(new_uri)
        ws.onopen = (event) => {
            if (address != "") {
                ws.send(JSON.stringify({"address": address, "action": "start"}))
            }
        }
        ws.onmessage = (event) => {
            console.log("messaging")
            console.log(event.data)
            setName(event.data)
        }
    }, [address])

    useEffect(() => {
        // connect to message websocket
        let new_uri = "ws://" + window.location.host + "/ws/chat/message"
        console.log(new_uri)
        let socket = new WebSocket(new_uri)
        socket.onmessage = (event) => {
            console.log(`message received ${event.data}`)
            setMessages([{text: event.data, fromMe:false}])
        }
        setSocket(socket)
    }, [])

    useEffect(() => {
        // connect to files websocket
        let new_uri = "ws://" + window.location.host + "/ws/chat/files"
        console.log(new_uri)
        let socket = new WebSocket(new_uri)
        socket.onmessage = (event) => {
            console.log(event, "!@$$%")
            // const file = axios.get(`/files/${event.data}`) // get file from storage
            setMessages([{text:"", isImage:true, img: "/files/" + event.data}])
            // setMessages([{text: event.data, fromMe:false}])
        }
        setFileSocket(socket)
    }, [])

    const onClickConnect = async (user: User) => {
        console.log(user.name)
        setIsChatLoading(true)
        const resp = await axios.post(`/api/v1/connect/${user.name}`)
        if (resp.status !== 200) {
            window.alert(`${resp.data.msg}`)
            return
        }
        setIsChatLoading(false)
        setAddress(resp.data["address"])
        setName(user.name)
        // setIsChatLoading(false)
        // setTimeout(() => {
        //     setConnected(true)
        //     setChosen(user)
        //     setIsChatLoading(false)
        // }, 5000)
    }

    const onClickEnd = (user: User) => {
        setTimeout(() => {
            setConnected(false)
            // setChosen({name: '', id: -1})
            setName("")
        }, 5000)
    }


    return isLoading ? (
        <Loader/>
    ) : (
        <div className="main-view-container">
            <div className="users-container">
                <div className="users-wrapper">
                    {users.map(user => {
                        return (
                            <UserComp user={user} onClickConnect={onClickConnect}
                                      onClickEnd={onClickEnd} connected={connected}/>
                        )
                    })}
                </div>
            </div>

            <div className="chat-container">
                {isChatLoading
                    ? <div className="chat-loader-wrapper">
                        <Loader />
                    </div>
                    : <div className="chat-wrapper">
                        <Chat name={name} msgs={messages} sock={msgSocket} fileSock={fileSocket}/>
                    </div>
                }
            </div>
        </div>
    )
}