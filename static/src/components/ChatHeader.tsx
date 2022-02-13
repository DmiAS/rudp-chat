import { AppBar, Box, IconButton, Toolbar, Typography } from '@material-ui/core'
import { report } from 'process'
import React from 'react'

export const ChatHead: React.FC = () => {
    return (
        <Box>
            <AppBar position="static">
                <Toolbar className="toolbar" variant="dense" style={{ background: '#FFDAB9' }}>
                    <Typography variant="h5" color="textPrimary" component="div">
                        Network
                    </Typography>
                </Toolbar>
            </AppBar>
        </Box>
    )
}
