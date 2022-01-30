import React from "react"
import {useContext} from "react";
import { useDispatch } from "react-redux";
import {Card, CardMedia, CardContent, CardActions, Typography, IconButton, Stack, Box} from "@mui/material"
import { AddShoppingCart } from "@mui/icons-material";

import { CartContext } from "../../../pages/App"
import { addCartItem } from "../../../actions";

const Product = ({product}) => {
    const dispatch = useDispatch()
    const value = useContext(CartContext)

    const handleCartClick = (productId, name, qty, price, desc, image) => {
        // Add +1 to qty if item already exists in cart
        let obj = value.cart.find(item => item.product_id === productId)

        if (obj !== undefined){
           qty += obj.qty
        }

        let item = {
            productId: productId,
            name: name,
            qty: qty,
            price: price,
            desc: desc,
            image: image,
        }

        dispatch(addCartItem(item)).then((result) => {
            value.onClickDrawer(true)
        })
    }

    return (
        <Card>
            <CardMedia image={product.image} title={product.name} style={{height:0, paddingTop: "56.29%"}}/>
            <CardContent>
                <Stack direction="row" spacing={2}>
                    <Typography variant="h6" gutterBottom>
                        {product.name}
                    </Typography>
                    <Box sx={{
                        flexGrow: 1
                    }}/>
                    <Typography variant="h6">
                        ${product.price}
                    </Typography>
                </Stack>
                <Typography variant="body2">{product.desc}</Typography>
            </CardContent>
            <CardActions>
                <Box sx={{
                    flexGrow: 1
                }}/>
                <IconButton aria-label="Add to cart" onClick={() => handleCartClick(
                    product.product_id,
                    product.name,
                    1,
                    product.price,
                    product.desc,
                    product.image)}>
                    <AddShoppingCart fontSize="large"/>
                </IconButton>
            </CardActions>
        </Card>
    )
}

export default React.memo(Product);