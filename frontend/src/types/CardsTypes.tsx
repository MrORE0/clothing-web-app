export type ItemType = {
    imgSrc: string,
    name: string,
    price: number
    // Probably will nead to add more properties later here too

}

export type CarouselItemType = {
    imgSrc: string,
    name: string,
    price: number
    // To know what preoperties will need to be added, ill have to compare 
    // the api return value for the products and the rendered values of each
    // individual brands cards for carousel wheels specifically

    // ***Add a link to the product as a nother property of the type***
}

export type SimplifiedItemType = {
    imgSrc: string,
    name: string,
    price: number
    // Probably will nead to add more properties later here too
}