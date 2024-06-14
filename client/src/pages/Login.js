import React, { useState, useEffect } from "react";
import axios from "axios";
import { useNavigate, Link } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import * as constants from "../utils/constant/Constant";

export default function Login() {
    const navigate = useNavigate();
    const toastOptions = {
        position: "top-right",
        autoClose: 3000,
        pauseOnHover: true,
        draggable: true,
        theme: "dark",
    };
    const [values, setValues] = useState({ email: "", password: "" });

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (token) {
            navigate("/");
        }
    }, [navigate]);

    const handleChange = (event) => {
        setValues({ ...values, [event.target.name]: event.target.value });
    };

    const validateForm = () => {
        const { email, password } = values;
        if (email === "" || password === "") {
            toast.error("Email and Password are required.", toastOptions);
            return false;
        }
        return true;
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (validateForm()) {
            const { email, password } = values;
            try {
                const { data } = await axios.post(`${process.env.URL}/signIn`, { email, password });
                console.log("----------->", data)
                localStorage.setItem("token", data.token);
                localStorage.setItem("userId", data.data.id);
                localStorage.setItem("name", data.data.name)
                navigate("/"); // Reload the page to trigger WebSocketProvider initialization
            } catch (error) {
                // console.log("----------->", data)
                console.error("Login Failed: ", error);
                toast.error(error, toastOptions);
            }
        }
    };


    return (
        <>
            <div className="login-body">
                <div className="form-container">
                    <form onSubmit={handleSubmit} className="form">
                        <h2 className="form-title">Log In</h2>
                        <input
                            type="text"
                            placeholder="Email"
                            autoComplete="off"
                            name="email"
                            onChange={handleChange}
                            className="form-input"
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            name="password"
                            onChange={handleChange}
                            className="form-input"
                        />
                        <button type="submit" className="form-button">
                            Log In
                        </button>
                        <p className="form-text">
                            Don't have an account?{" "}
                            <Link to="/signup" className="form-link">
                                Create One
                            </Link>
                        </p>
                    </form>
                </div>
            </div>
            <ToastContainer />
        </>
    );
}
