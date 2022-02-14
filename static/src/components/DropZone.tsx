import React from 'react'
import { useDropzone } from 'react-dropzone'

export const DropZone: React.FC = () => {
    const onDrop = React.useCallback((acceptedFiles) => {
        console.log(acceptedFiles)
    }, [])



    const { getRootProps, getInputProps, isDragAccept, isDragActive } = useDropzone({
        onDrop,
        multiple: false,
        accept: "image/jpeg,image/png"
    })

    return (
        <div className={isDragActive ? "dropzone-wrapper active" : "dropzone-wrapper"}>
            <div {...getRootProps()} className="dropzone">
                <input {...getInputProps()} />
                <p className='dd-text'>Drag & Drop Files Here</p>
            </div>
        </div>
    )
}