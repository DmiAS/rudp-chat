import { AppBar, Box, IconButton, Toolbar, Typography } from "@material-ui/core";
import CloseIcon from '@material-ui/icons/Close';
import React from "react";
import { useNavigate } from "react-router-dom";

export const ToolBar: React.FC = () => {
    const navigate = useNavigate()

    return (<Box>
        <AppBar position="static">
            <Toolbar className="toolbar" variant="dense" style={{background: '#FFDAB9' }}>
                <Typography variant="h5" color="textPrimary" component="div">
                    Network
                </Typography>
                <IconButton onClick={() => navigate('/')} edge="end" size="medium" style={{ color: "#ff0000" }}>
                    <CloseIcon />
                </IconButton>
            </Toolbar>
        </AppBar>
    </Box>)
}
