import React from "react";
import {Card, CardMedia, CardContent, CardActions, Typography, IconButton} from "@mui/material"

const Product = ({product}) => {
    return (
        <Card>
            {console.log(product.image)}
            <CardMedia image={product.image} title={product.name} style={{height:0, paddingTop: "56.25%"}}/>
            <CardContent>
                <div>
                    <Typography variant="h5" gutterBottom>
                        {product.name}
                    </Typography>
                    <Typography variant="h5">
                        {product.price}
                    </Typography>
                    <Typography variant="h5">
                        {product.qty}
                    </Typography>
                </div>
            </CardContent>
            <CardActions>
                <IconButton aria-label="Add to cart">
                    
                </IconButton>
            </CardActions>
        </Card>
    )
}

export default Product;