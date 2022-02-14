import { AppBar, Box } from '@material-ui/core'
import { report } from 'process'
import React from 'react'
import { Msg, User } from '../interfaces/core'
import { ChatFooter } from './ChatFooter'
import { ChatHeader } from './ChatHeader'
import { ChatMsg } from './ChatMsg'

type Props = {
    user: User
    msgs: Msg[]
}

export const Chat: React.FC<Props> = ({ user, msgs }) => {
    return (
        <>
            <div className="chat-header">
                <ChatHeader user={user}/>
            </div>

            <div className="chat-msg">
                <ChatMsg msgs={msgs}/>
            </div>

            <div className="chat-footer">
                <ChatFooter />
            </div>
        </>
    )
}
