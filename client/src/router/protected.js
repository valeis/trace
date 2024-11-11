import {Navigate, Outlet} from "react-router-dom";

const Protected = () => {
    const token = localStorage.getItem("access-token");

    return token ? <Outlet /> : <Navigate to= "/login" />;
};

export default Protected;