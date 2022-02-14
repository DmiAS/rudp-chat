import React, { useEffect } from "react"
import { Msg, User } from "../interfaces/core"
import { Loader } from "./Loader"
import { UserComp } from "./UserComp"
import axios from "axios"
import { Chat } from "./Chat"

type ConnectionProps = {
    users: User[]
    isLoading: boolean
}

const msgs: Msg[] = [{text: "Привет"}, {text: "1234123412341234123412341234", fromMe: true}]

export const Connection: React.FC<ConnectionProps> = ({ users, isLoading }) => {
    const [connected, setConnected] = React.useState(false)
    const [isChatLoading, setIsChatLoading] = React.useState(false)
    const [chosen, setChosen] = React.useState({ name: '', id: -1 })
    const [address, setAddress] = React.useState("")

    useEffect(() => {
        let new_uri = "ws://" + window.location.host + "/ws/chat/thread"
        console.log(new_uri)
        let ws = new WebSocket(new_uri)
        ws.onopen = (event) => {
            ws.send(JSON.stringify({ "address": address, "action": "start" }))
        }
    }, [address])

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
            setChosen({ name: '', id: -1 })
        }, 2000)
    }


    return isLoading ? (
        <Loader />
    ) : (
        <div className="main-view-container">
            <div className="users-container">
                <div className="users-wrapper">
                    {users.map(user => {
                        return (
                            <UserComp user={user} onClickConnect={onClickConnect}
                                onClickEnd={onClickEnd} chosen={chosen} connected={connected} />
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
                        <Chat user={chosen} msgs={msgs}/>
                    </div>
                }
            </div>
        </div>
    )
}