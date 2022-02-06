import React from "react"
import {Box, Typography} from "@mui/material";
import { Link } from "react-router-dom";

import { BlackButton } from "../../../theme";

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
                <BlackButton variant="outlined" component={Link} to="/" text="Back"/>
            </Box>
        </>
    )
}

export default Confirmation