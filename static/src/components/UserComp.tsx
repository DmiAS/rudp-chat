import { Accordion, AccordionDetails, AccordionSummary, Avatar, Button, Typography } from '@material-ui/core'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import React, { useEffect } from 'react'
import { User } from '../interfaces/core'

interface UserActions {
    user: User
    onClickConnect: (user: User) => void
    onClickEnd: (user: User) => void
    chosen: User
    connected: boolean
}

export const UserComp: React.FC<UserActions> = ({ user, onClickConnect, onClickEnd, chosen, connected }) => {

    React.useEffect(() => {
        console.log(user.id !== chosen.id, chosen)
    })

    return (
        <Accordion key={user.id} className='user-record'>
            <AccordionSummary
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel1a-content"
                id="panel1a-header"
            >
                <Typography>
                    <div className='avatar-name'>
                        <Avatar />
                        <div style={{marginLeft: "30px"}}>{user.name}</div>
                    </div>
                </Typography>
            </AccordionSummary>
            <AccordionDetails className='actions-container'>
                <div className="actions">
                    <Button disabled={connected} onClick={() => onClickConnect(user)} color="primary" variant="outlined" size="small">Connect</Button>
                    <Button disabled={!connected || chosen.id !== user.id}  onClick={() => onClickEnd(user)} style={{ color: "#ff0000" }} variant="outlined" size="small">End</Button>
                </div>
            </AccordionDetails>
        </Accordion>
    )
}