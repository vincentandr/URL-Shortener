import React from "react";

import { Cart, Login} from "../../components";
import Bar from "./Bar/Bar";

const Navbar = ({cart, drawer, login}) => {
    return (
        <>
            <Bar cart={cart} drawer={drawer} login={login}/>
            <Cart cart={cart} drawer={drawer}/>
            <Login login={login}/>
        </>
    )
}

export default React.memo(Navbar)