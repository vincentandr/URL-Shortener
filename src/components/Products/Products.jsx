import React from "react"
import {useState, useEffect} from "react";
import {Grid, Snackbar} from "@mui/material";

import Product from "./Product";

const Products = ({products}) => {
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

    return (
        <main>
            <Snackbar
                open={open}
                autoHideDuration={6000}
                message="Item added to cart"
                TransitionProps={{onExited: handleExited}}
                severity="success"
                />
            <Grid container justify="center" spacing={4}>
                {products.map((product) => (
                    <Grid item key={product.product_id} xs={6} md={4} lg={3}>
                        <Product product={product} setSnackPack={setSnackPack}/>
                    </Grid>
                ))}
            </Grid>
        </main>
    )
}

export default React.memo(Products);