import { AppBar, Box, IconButton, TextField, Toolbar, Typography } from '@material-ui/core'
import { report } from 'process'
import React from 'react'
import { User } from '../interfaces/core'
import InputBase from '@material-ui/icons/Input';
import MenuIcon from '@material-ui/icons/Menu';
import SearchIcon from '@material-ui/icons/Search';
import { Search } from '@material-ui/icons';
import SendIcon from '@material-ui/icons/Send';


export const ChatFooter: React.FC = () => {
    return (
        <Box>
            <AppBar position="static">
                <Toolbar style={{ background: "#1976d2", minHeight: "100px" }}>
                    <Typography style={{width: "100%"}}>
                        <div className="chat-input">
                            <TextField label="Type message" variant="filled" style={{ background: "#FFFFFF4A", marginRight: "30px", flexGrow: "7" }}/>
                            <IconButton>
                                <SendIcon style={{color: "#FFF"}}/>
                            </IconButton>
                        </div>
                    </Typography>

                </Toolbar>
            </AppBar>
        </Box>
    )
}
