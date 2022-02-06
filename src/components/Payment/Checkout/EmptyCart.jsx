import React from "react"
import { Box, Stack, Typography } from "@mui/material"
import { Link } from "react-router-dom";
import { Logo } from "../../../theme";

const EmptyCart = () => {
    return (
        <Stack spacing={2} justifyContent="center" alignItems="center" textAlign="center" height="100vh" sx={{
            p:2
        }}>
            <div>
                <Logo component={Link} to="/" gutterBottom variant="h3" color="inherit"/>
            </div>
            <Box 
                component="img"
                sx={{
                    maxHeight: { xs: 250, md: 275 },
                    maxWidth: { xs: 250, md: 275 },
                }}
                alt="product img"
                src="https://www.metro-markets.com/plugins/user/images/emptycart.png"
            />
            <Typography variant="h6" gutterBottom>
                You don't have any item in your shopping cart. 
                <Box component={Link} to="/" sx={{
                    textDecoration: "none",
                }}> Click here </Box>
                to start shopping now!
            </Typography>
        </Stack>
    )
}

export default EmptyCart