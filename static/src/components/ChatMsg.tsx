import { Button } from "@material-ui/core"
import React from "react"
import { Msg } from "../interfaces/core"

type Props = {
    msgs: Msg[]
}

export const ChatMsg: React.FC<Props> = ({ msgs }) => {
    console.log(msgs)

    return (
        <div>
            {msgs.map(elem => {
                return (
                    <div className={elem.fromMe ? "msg reverse" : "msg"}>
                        {/* <div className={elem.fromMe ? "msg-cloud green" : "msg-cloud blue"}>
                            {elem.text}
                        </div> */}
                        <Button style={{minWidth: "160px"}} size="large" variant="outlined" color={elem.fromMe ? "primary" : "inherit"}>{elem.text}</Button>
                    </div>
                )
            })}
        </div>
    )
}