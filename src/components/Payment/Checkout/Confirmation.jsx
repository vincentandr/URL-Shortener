import React from "react"
import {Box, Typography, Button} from "@mui/material";
import { Link } from "react-router-dom";

const Confirmation = ({payment}) => {
    return (
        <>
            <Typography variant="h5">
                Order Confirmed!
            </Typography>
            <Box sx={{
                display:"flex",
                justifyContent: "end",
                pt: 3,
                pb: 3,
            }}>
                <Button variant="outlined" component={Link} to="/">Back</Button>
            </Box>
        </>
    )
}

export default Confirmation