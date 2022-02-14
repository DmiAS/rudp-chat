import { Button } from "@material-ui/core"
import React from "react"
import { Msg } from "../interfaces/core"

type Props = {
    msgs: Msg[]
}

export const ChatMsg: React.FC<Props> = ({ msgs }) => {
    console.log(msgs)

    const addClass = (elem: Msg) => {
        let res: string = ''

        if (elem.fromMe) {
            res += "msg reverse"
        } else {
            res += "msg"
        }

        if (elem.isImage) {
            res += " image"
        }

        return res
    }

    return (
        <>
            {msgs.map(elem => {
                return (
                    <div className={addClass(elem)}>
                        {!elem.isImage ?
                            <Button
                                size="large"
                                variant="outlined"
                                color={elem.fromMe ? "primary" : "inherit"}
                                style={{minWidth: "160px"}}
                            >
                                {elem.text}
                            </Button>
                        :
                            <img src={elem.img} style={{width: "520px", height: "320px"}} />
                        }

                    </div>
                )
            })}
        </>
    )
}