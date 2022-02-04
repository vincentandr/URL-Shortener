import React from "react"
import {Typography, Button} from "@mui/material";
import { Link } from "react-router-dom";

const Confirmation = ({cart}) => {
    return (
        <>
            <Typography variant="h5">
                Order Confirmed
            </Typography>
            <Button variant="outlined" component={Link} to="/">Back</Button>
        </>
    )
}

export default Confirmation