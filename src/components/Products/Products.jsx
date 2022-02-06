import React from "react"
import {useState, useEffect} from "react";
import {Grid, Snackbar, Alert, Box} from "@mui/material";

import Product from "./Product/Product";

const Products = ({products, cart, onClickDrawer}) => {
    const [snackPack, setSnackPack] = useState([]);
    const [open, setOpen] = useState(false);
    const [messageInfo, setMessageInfo] = useState(undefined);

    useEffect(() => {
        if (snackPack.length && !messageInfo) {
        // Set a new snack when we don't have an active one
        setMessageInfo({ ...snackPack[0] });
        setSnackPack((prev) => prev.slice(1));
        setOpen(true);
        } else if (snackPack.length && messageInfo && open) {
        // Close an active snack when a new one is added
        setOpen(false);
        }
    }, [snackPack, messageInfo, open]);

    const handleExited = () => {
        setMessageInfo(undefined)
    }

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
        return;
        }

        setOpen(false);
    };

    return (
        <main>
            <Snackbar
                open={open}
                autoHideDuration={6000}
                TransitionProps={{onExited: handleExited}}
                severity="success"
                onClose={handleClose}
            >
                <Alert severity="success" sx={{ width: '100%' }}>
                    Item added to cart
                </Alert>
            </Snackbar>
            <Box sx={{
                p: 3
            }}>
                <Grid container justify="center" spacing={4}>
                    {products.map((product) => (
                        <Grid item key={product.product_id} xs={12} sm={6} lg={3}>
                            <Product product={product} cart={cart} onClickDrawer={onClickDrawer} setSnackPack={setSnackPack}/>
                        </Grid>
                    ))}
                </Grid>
            </Box>
        </main>
    )
}

export default React.memo(Products);