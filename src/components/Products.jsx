import React, {useState, useEffect} from "react";
import {Grid} from "@mui/material";
import axios from "axios";

import Product from "./Product";

const Products = () => {
    const [products, setProducts] = useState([])

    const uri = "/products"

    useEffect(() => {
        axios.get(uri).then(function (response) {
            if (response.status == 200) {
                setProducts(response.data.products)
            }
        })
    }, [])

    return (
        <main>
            <h2>test</h2>
            <Grid container justify="center" spacing={4}>
                {products.map((product) => (
                    <Grid item key={product.product_id} xs={6} md={4} lg={3}>
                        <Product product={product}/>
                    </Grid>
                ))}
            </Grid>
        </main>
    )
}

export default Products;