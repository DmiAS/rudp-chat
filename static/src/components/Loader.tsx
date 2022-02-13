import { CircularProgress } from "@material-ui/core"
import React from "react"


export const Loader: React.FC = () => {
    return (
        <div className="loader-wrapper">
            <CircularProgress size={100} />
        </div>
    )
}