import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { Connection } from './components/Connection';
import { StartBtn } from './components/StartBtn';
import { ToolBar } from './components/ToolBar';
import { User } from './interfaces/core';

function App() {
    const [users, setUsers] = React.useState<User[]>([])
    const [isLoading, setLoading] = React.useState(true)

    const namesHandler = (names: string[]) => {
        const newMas = names.map(elem => {
            return {name: elem, id: Date.now()}
        })

        setUsers(prev => [...newMas])
    }

    React.useEffect(() => {

        setLoading(false)
    }, [users])

    return (
        <Router>
            <div className='main-container'>
                <ToolBar />

                <div className='main-content'>
                    <Routes>
                        <Route element={<StartBtn namesHandler={namesHandler}/>} path="/" />
                        <Route element={<Connection isLoading={isLoading} users={users} />} path="/connect" />
                    </Routes>
                </div>
            </div>

        </Router>
    )
}

export default App;
