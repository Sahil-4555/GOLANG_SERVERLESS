import React from "react";


export default function Home() {
    const userid = localStorage.getItem("userId")
    return (
        <div>
            {userid}
        </div>
    )
}