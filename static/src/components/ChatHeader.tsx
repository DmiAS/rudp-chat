import { AppBar, Box, IconButton, Toolbar, Typography } from '@material-ui/core'
import { report } from 'process'
import React from 'react'
import { User } from '../interfaces/core'

type Props = {
    name: string
}

export const ChatHeader: React.FC<Props> = ({ name }) => {

    return (
        <Box>
            <AppBar position="static">
                <Toolbar className="toolbar" style={{background: "#1976d2", minHeight: "80px"}}>
                    <Typography variant="h2" color="textPrimary" component="div">
                        <div style={{color: "#FFFFF0"}}>
                            Dialog to {name}
                        </div>
                    </Typography>
                </Toolbar>
            </AppBar>
        </Box>
    )
}
