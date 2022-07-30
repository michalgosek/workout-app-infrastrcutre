import { Footer, NavBar } from "../components";
import React, { PropsWithChildren } from "react";

interface PageLayoutProps {
    children: React.ReactNode;
}

const PageLayout: React.FC<PropsWithChildren<PageLayoutProps>> = ({ children }) => {
    return (
        <div id="app" className="d-flex flex-column h-100">
            <NavBar />
            <div className="container flex-grow-1">
                <div className="mt-5">
                    {children}
                </div>
            </div>
            <Footer />
        </div>
    );
};

export default PageLayout;