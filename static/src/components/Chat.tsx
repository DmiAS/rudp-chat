import React from 'react'
import {Msg, User} from '../interfaces/core'
import {ChatFooter} from './ChatFooter'
import {ChatHeader} from './ChatHeader'
import {ChatMsg} from './ChatMsg'

type Props = {
    name: string
    msgs: Msg[]
    sock?:WebSocket
}

export const Chat: React.FC<Props> = ({name, msgs, sock}) => {
    const [messages, setMessages] = React.useState<Msg[]>([])

    const msgHandler = (msg: Msg) => {
        sock!.send(msg.text)
        setMessages(prev => [...prev, msg])
    }

    React.useEffect(() => {
        // console.log(`in use ref with messages = ${msgs}`)
        setMessages(prev => [...prev, ...msgs])
    }, [msgs])

    return (
        <>
            <div className="chat-header">
                <ChatHeader name={name}/>
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
