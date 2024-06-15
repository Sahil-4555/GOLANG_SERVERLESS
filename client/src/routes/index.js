import { Routes, Route, BrowserRouter } from 'react-router-dom';
import { lazy } from 'react';
import Middleware from './middleware';
import React from 'react';

const SomethingwentWrong = lazy(() => import('../pages/SomethingwentWrong'));
const Login = lazy(() => import('../pages/Login'));
const Signup = lazy(() => import('../pages/Signup'));
const TodoWrapper = lazy(() => import('../pages/TodoWrapper'))

const RouteList = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Middleware />}>
                    <Route path='/' exact element={<TodoWrapper />} />
                </Route>
                <Route exact path="/error" element={<SomethingwentWrong />} />
                <Route exact path="/login" element={<Login />} />
                <Route exact path="/signup" element={<Signup />} />
            </Routes>
        </BrowserRouter>
    );
};

export default RouteList