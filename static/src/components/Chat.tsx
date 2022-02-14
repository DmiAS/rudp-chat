import React from 'react'
import {Msg, User} from '../interfaces/core'
import {ChatFooter} from './ChatFooter'
import {ChatHeader} from './ChatHeader'
import {ChatMsg} from './ChatMsg'
import { DropZone } from './DropZone'

type Props = {
    name: string
    msgs: Msg[]
    sock?:WebSocket
    fileSock?: WebSocket
}

export const Chat: React.FC<Props> = ({name, msgs, sock, fileSock}) => {
    const [messages, setMessages] = React.useState<Msg[]>([])
    const [images, setImages] = React.useState<any[]>([])
    const [imageURLs, setImageURLs] = React.useState<any[]>([])

    const msgHandler = (msg: Msg) => {
        sock!.send(msg.text)
        setMessages(prev => [...prev, msg])
    }

    const uploadFile = (event: any) => {
        let file = event.target.files[0]
        console.log(file)

        if (file) {
            setImages([file])

            let reader = new FileReader();
            reader.readAsDataURL(file);
            reader.onload = function () {
            //me.modelvalue = reader.result;
                console.log(reader.result);
                const postJSON = {name: file.name, data: reader.result!.slice(String(reader.result).indexOf(',') + 1)}
                console.log(postJSON.data)
                fileSock!.send(JSON.stringify(postJSON))
            };
            reader.onerror = function (error) {
                console.log('Error: ', error);
            };
        }
    }

    React.useEffect(() => {
        if (images.length < 1) return;

        const newImageUrls: any[] = []
        images.forEach(image => newImageUrls.push(URL.createObjectURL(image)))
        setImageURLs(newImageUrls)
    }, [images])

    React.useEffect(() => {
        setMessages(prev => [...prev, ...imageURLs.map(elem => {
            return {text: '', isImage: true, img: elem, fromMe: true}
        })])
    }, [imageURLs])

    //       {imageURLs.map(imgScr => <img src={imgScr} style={{width: "280px"}} />)}

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
                <ChatFooter uploadFile={uploadFile} msgHandler={msgHandler}/>
            </div>

            {/* <div className="drop-zone">
                <DropZone />
            </div> */}
        </>
    )
}
