import React from "react";
import {Container, Typography, Button, Drawer, Stack, useMediaQuery, useTheme} from "@mui/material"
import {Close} from "@mui/icons-material"

import EmptyCart from "./EmptyCart/EmptyCart";
import FilledCart from "./FilledCart/FilledCart";

const CartContent = ({cart, drawer}) => {
    const isEmpty = !cart.products.length
    
    return (
        <Container sx={{
            height: "100%",
        }}>
            <Stack direction="row" spacing={2}>
                <Typography variant="h4" width="100%">
                    Shopping Cart
                </Typography>
                <Button variant="text" color="inherit" size="large" startIcon={<Close />} onClick={() => drawer.onClick(false)}>
                    Close
                </Button>
            </Stack>
            {isEmpty ? <EmptyCart/> : <FilledCart cart={cart} drawer={drawer}/>}
        </Container>
    )
}

const Cart = ({cart, drawer}) => {
    const theme = useTheme()
    const largeScreen = useMediaQuery(theme.breakpoints.up("sm"))

    return(
        <Drawer
            anchor="right"
            open={drawer.state}
            onClose={() => drawer.onClick(false)}
            PaperProps={largeScreen ? {
                sx: {
                    width: 450,
                }
            } : {
                sx: {
                    width: "100%",
                }
            }
            }
        >
            <CartContent cart={cart} drawer={drawer}/>
        </Drawer>
    )
}

export default React.memo(Cart)