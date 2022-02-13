import React from "react"
import { User } from "../interfaces/core"
import { Chat } from "./Chat"
import { Loader } from "./Loader"
import { UserComp } from "./UserComp"

type ConnectionProps = {
    users: User[]
    isLoading: boolean
}


export const Connection: React.FC<ConnectionProps> = ({ users, isLoading }) => {
    const [connected, setConnected] = React.useState(false)
    const [isChatLoading, setIsChatLoading] = React.useState(false)
    const [chosen, setChosen] = React.useState({name: '', id: -1})

    const onClickConnect = (user: User) => {
        setIsChatLoading(true)
        setTimeout(() => {
            setConnected(true)
            setChosen(user)
            setIsChatLoading(false)
        }, 5000)
    }

    const onClickEnd = (user: User) => {
        setTimeout(() => {
            setConnected(false)
            setChosen({name: '', id: -1})
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
                        <Chat />
                    </div>
                }
            </div>
        </div>
    )
}