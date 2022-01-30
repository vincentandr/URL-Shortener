import React from "react";
import {Toolbar, IconButton, Badge, Box, Typography, AppBar, TextField, InputAdornment, Stack} from "@mui/material";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import {ShoppingCart, Search, AccountCircle} from "@mui/icons-material";
import { useDispatch } from "react-redux";

import { searchProducts } from "../../actions/Products";

const theme = createTheme({
  typography: {
    fontFamily: 'Fugaz One',
  },
});


const Navbar = ({totalItems, onClickDrawer}) => {
    const dispatch = useDispatch();

    const search = (e) => {
        dispatch(searchProducts(e.target.value))
    }

    return (
        <AppBar position="relative" color="inherit">
            <Toolbar sx={{
                ml: "5%",
                mr: "5%",
            }}>
                <ThemeProvider theme={theme}>
                    <Typography variant="h4" color="inherit">
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
                    <IconButton aria-label="Show cart items" color="inherit" onClick={() => onClickDrawer(true)}>
                        <Badge badgeContent={totalItems} color="secondary">
                            <ShoppingCart/>
                        </Badge>
                    </IconButton>
                    <IconButton aria-label="Show cart items"  color="inherit">
                        <AccountCircle/>
                    </IconButton>
                </Stack>
        </Toolbar>
        </AppBar>
    )
}

export default React.memo(Navbar)