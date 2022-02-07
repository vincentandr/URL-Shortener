import React from "react";
import { useDispatch } from "react-redux";
import { Link } from "react-router-dom";
import {Toolbar, IconButton, Badge, Box, AppBar, TextField, InputAdornment, Stack, useMediaQuery, useTheme} from "@mui/material";
import {ShoppingCart, Search, AccountCircle} from "@mui/icons-material";

import { searchProducts } from "../../../actions/Products";
import { Logo } from "../../../theme";

const Bar = ({cart, drawer, login}) => {
    const dispatch = useDispatch();
    const theme = useTheme()
    const smallPhone = useMediaQuery(theme.breakpoints.down("sm"))

    const search = (e) => {
        dispatch(searchProducts(e.target.value))
    }

    return(
        <AppBar position="sticky" color="inherit">
            <Toolbar>
                <Logo component={Link} to="/" variant="h4" color="inherit"/>
                <Box sx={{
                    flexGrow: 1
                }}/>
                <Stack direction="row" spacing={1}>
                    {!smallPhone && <TextField
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
                    />}
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
            { smallPhone &&
            <Box sx={{
                pl: 3,
                pr: 3,
                pb: 2,
                pt: 1,
            }}>
                <TextField
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <Search />
                            </InputAdornment>
                        ),
                    }}
                    fullWidth
                    placeholder="Search..."
                    onChange={search}
                    size="small"
                    variant="outlined"
                />
            </Box>
        }
        </AppBar>
    )
}

export default Bar