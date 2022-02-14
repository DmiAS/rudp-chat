import { AppBar, Box, IconButton, TextField, Toolbar, Typography } from '@material-ui/core'
import { report } from 'process'
import React from 'react'
import { Msg, User } from '../interfaces/core'
import InputBase from '@material-ui/icons/Input';
import MenuIcon from '@material-ui/icons/Menu';
import SearchIcon from '@material-ui/icons/Search';
import { Search } from '@material-ui/icons';
import SendIcon from '@material-ui/icons/Send';
import FileUploadIcon from '@material-ui/icons/CloudUpload';

type Props = {
    msgHandler: (text: Msg) => void
    uploadFile: (event: any) => void
}

export const ChatFooter: React.FC<Props> = ({ msgHandler, uploadFile }) => {
    // const [images, setImages] = React.useState<any[]>([])


    // React.useEffect(() => {
    //     if (images.length < 1) return;

    //     const newImageUrls: any[] = []
    //     images.forEach(image => newImageUrls.push(URL.createObjectURL(image)))
    //     setImageURLs(newImageUrls)
    // }, [images])

    const ref = React.useRef<HTMLInputElement>(null)
    const refOpen = React.useRef<HTMLInputElement>(null)

    const handler = () => {
        const text = ref.current!.value
        msgHandler({text: text, fromMe: true})
        ref.current!.value = ''
    }

    const uploadHandler = () => {
        refOpen.current!.click();
    }

    // const uploadFile = (event: any) => {
    //     let file = event.target.files[0]
    //     console.log(file)

    //     if (file) {
    //         setImages([file])
    //     }
    // }

    return (
        <Box>
            <AppBar position="static">
                <Toolbar style={{ background: "#1976d2", minHeight: "100px" }}>
                    <Typography style={{width: "100%"}}>
                        <div className="chat-input">
                            <TextField inputRef={ref} label="Type message" variant="filled" style={{ background: "#FFFFFF4A", marginRight: "30px", flexGrow: "7" }}/>
                            <IconButton onClick={uploadHandler}>
                                <FileUploadIcon style={{color: "#FFF"}}/>
                            </IconButton>
                            <IconButton onClick={handler}>
                                <SendIcon style={{color: "#FFF"}}/>
                            </IconButton>
                        </div>
                    </Typography>

                </Toolbar>
            </AppBar>

            <input type='file' id='file' ref={refOpen} onChange={e => uploadFile(e)} style={{display: 'none'}}/>
        </Box>
    )
}
