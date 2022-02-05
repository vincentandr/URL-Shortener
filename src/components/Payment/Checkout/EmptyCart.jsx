import React from "react"
import { Box, Paper, Typography } from "@mui/material"
import {ThemeProvider, createTheme } from "@mui/material/styles";
import { Link } from "react-router-dom";

const theme = createTheme({
    typography: {
        fontFamily: 'Fugaz One',
    },
});

const EmptyCart = () => {
    return (
        <Box sx={{
            display: "flex",
            justifyContent: "center",
            marginTop: 5,
        }}>
            <Paper sx={{
                width: 1/3,
                pl: 3,
                pr: 3,
                textAlign: "center",}}>
                    <ThemeProvider theme={theme}>
                        <Typography component={Link} to="/" gutterBottom variant="h4" color="inherit" sx={{
                            textDecoration: "none",
                        }}>
                            Microshopping
                        </Typography>
                    </ThemeProvider>
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
            </Paper>
        </Box>
    )
}

export default EmptyCart