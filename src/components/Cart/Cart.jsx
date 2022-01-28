import React from "react";
import {useState} from "react";
import {Box, Container, Typography, Button, Table, TableContainer, TableBody, TableRow, TableCell, Paper} from "@mui/material"

const Cart = ({cart}) => {
    const isEmpty = !cart.length

    const EmptyCart = () => (
        <Typography variant="h5">
            You have no item in your shopping cart. Begin your shopping experience now!
        </Typography>
    )

    const FilledCart = () => (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableBody>
                {cart.map((product) => (
                    <TableRow
                    key={product.product_id}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                    >
                    <TableCell aling="center">
                        <Box component="img" sx={{
                            height: 233,
                            width: 350,
                            maxHeight: { xs: 233, md: 167 },
                            maxWidth: { xs: 350, md: 250 },
                            }}
                            alt="The house from the offer."
                            src={product.image}/>
                    </TableCell>
                    <TableCell align="right">{product.name}</TableCell>
                    <TableCell align="right">{product.qty}</TableCell>
                    <TableCell align="right">{product.price}</TableCell>
                    <TableCell align="right">{product.desc}</TableCell>
                    </TableRow>
                ))}
                </TableBody>
            </Table>
        </TableContainer>
    )

    return(
        <Container>
            <Typography variant="h2">
                Shopping Cart
            </Typography>
            {isEmpty ? <EmptyCart/> : <FilledCart/>}
            <div>
                <Button type="button" variant="contained" color="secondary">Empty Cart</Button>
                <Button type="button" variant="contained" color="primary">Checkout</Button>
            </div>
        </Container>
    )
}

export default Cart