import React from "react";
import { useDispatch } from "react-redux";
import { Link } from "react-router-dom";
import {Toolbar, IconButton, Badge, Box, Typography, AppBar, TextField, InputAdornment, Stack} from "@mui/material";
import {ThemeProvider, createTheme } from "@mui/material/styles";
import {ShoppingCart, Search, AccountCircle} from "@mui/icons-material";

import { searchProducts } from "../../../actions/Products";

const Bar = ({cart, drawer, login}) => {
    const dispatch = useDispatch();

    const search = (e) => {
        dispatch(searchProducts(e.target.value))
    }

    const theme = createTheme({
        typography: {
            fontFamily: 'Fugaz One',
        },
    });

    return(
        <AppBar position="relative" color="inherit">
            <Toolbar sx={{
                ml: "5%",
                mr: "5%",
            }}>
                <ThemeProvider theme={theme}>
                    <Typography component={Link} to="/" variant="h4" color="inherit" sx={{
                        textDecoration: "none"
                    }}>
                        Microshopping
                    </Typography>
                </ThemeProvider>
                <Box sx={{
                    flexGrow: 1
                }}/>
                <Stack direction="row" spacing={1}>
                    <TextField
                        id="outlined-adornment-password"
                        InputProps={{
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Search />
                                </InputAdornment>
                            ),
                        }}
                        placeholder="Search..."
                        onChange={search}
                        size="small"
                        variant="outlined"
                    />
                    <IconButton aria-label="Show cart items" color="inherit" onClick={() => drawer.onClick(true)}>
                        <Badge badgeContent={cart.products.length} color="secondary">
                            <ShoppingCart/>
                        </Badge>
                    </IconButton>
                    <IconButton aria-label="Show cart items"  color="inherit" onClick={() => login.onClick(true)}>
                        <AccountCircle/>
                    </IconButton>
                </Stack>
        </Toolbar>
        </AppBar>
    )
}

export default Bar