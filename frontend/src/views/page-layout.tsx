import { Footer, NavBar } from "../components";
import { Outlet } from "react-router-dom";
import React from "react";

const PageLayout: React.FC = () => {
    return (
        <div id="app" className="d-flex flex-column h-100">
            <NavBar />
            <div className="container flex-grow-1">
                <div className="mt-5">
                    < Outlet />
                </div>
            </div>
            <Footer />
        </div >
    );
};

export default PageLayout;
