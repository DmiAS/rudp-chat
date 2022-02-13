import { Button, TextField } from "@material-ui/core"
import axios from "axios"
import React from "react"
import { useNavigate } from "react-router-dom"

interface ApiResult {
    names: string[]
}

interface Props {
    namesHandler: (names: string[]) => void
}

export const StartBtn: React.FC<Props> = props => {
    const ref = React.useRef<HTMLInputElement>(null)

    const textHandler = async () => {
        const value = ref.current!.value
        const result = await axios.post(`http://localhost:8080/api/v1/register/${value}`)
        if (result.status !== 200) {
            window.alert(`${result.data.msg}`)
        }

        // !!! Зачекать потом
        const jsonResult = await axios.get<ApiResult>(`http://localhost:8081/api/v1/users/${value}`)
        const names = jsonResult.data.names
        console.log(names)


        props.namesHandler(names)

        props.namesHandler(["Dima", "masha", "pasha"])

        navigate('/connect')
    }

    const navigate = useNavigate()
    return (
        <div className="btn-wrapper">
            <TextField inputRef={ref} style={{marginBottom: "20px"}} id="filled-basic" label="Filled" variant="filled" />

            <Button onClick={textHandler} size='large' color='primary' variant='contained' className='initial-btn' >
                Start
            </Button>
        </div>
    )
}