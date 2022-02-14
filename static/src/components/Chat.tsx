import { AppBar, Box } from '@material-ui/core'
import { report } from 'process'
import React, { useEffect } from 'react'
import { Msg, User } from '../interfaces/core'
import { ChatFooter } from './ChatFooter'
import { ChatHeader } from './ChatHeader'
import { ChatMsg } from './ChatMsg'

type Props = {
    user: User
    msgs: Msg[]
}

export const Chat: React.FC<Props> = ({ user, msgs }) => {
    const [messages, setMessages] = React.useState<Msg[]>([])

    const msgHandler = (msg: Msg) => {
        setMessages(prev => [...prev, msg])
    }

    React.useEffect(() => {
        setMessages(prev => [...prev, ...msgs])
    }, [])

    return (
        <>
            <div className="chat-header">
                <ChatHeader user={user}/>
            </div>

            <div className="chat-msg">
                <ChatMsg msgs={messages}/>
            </div>

            <div className="chat-footer">
                <ChatFooter msgHandler={msgHandler}/>
            </div>
        </>
    )
}
