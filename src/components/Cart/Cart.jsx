import React from "react";
import {Container, Typography, Button, Drawer, Stack} from "@mui/material"
import {Close} from "@mui/icons-material"

import EmptyCart from "./EmptyCart/EmptyCart";
import FilledCart from "./FilledCart/FilledCart";

const Cart = ({cart, drawer}) => {
    const isEmpty = !cart.products.length

    return(
        <Drawer
            anchor="right"
            open={drawer.state}
            onClose={() => drawer.onClick(false)}
            PaperProps={{
                sx: { width: "30%" },
            }}
        >
            <Container>
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
        </Drawer>
    )
}

export default React.memo(Cart)