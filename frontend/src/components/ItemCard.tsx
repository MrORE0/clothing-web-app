import type { ItemType, CarouselItemType, SimplifiedItemType } from "../types/CardsTypes"

//Interfaces for card components to structure the passed dynamic data
interface ProductCardInt {
    itemData: ItemType;
}

interface CarouselCardInt {
    itemData: CarouselItemType;
}

interface SimlpifiedCardInt {
    itemData: SimplifiedItemType;
}

export const ProductCard: React.FC<ProductCardInt> = (itemData) => {
    return (
        <>
            {/* Basic card, will need to match the cards from the brands websites data */}
        </>
    )

}

export const CarouselCard: React.FC<CarouselCardInt> = (itemData) => {
    return (
        <>
            {/* ***Need to convert the carousel card code here***  */}
        </>
    )

}

export const SimplifiedCard: React.FC<SimlpifiedCardInt> = (itemData) => {
    return (
        <>
            {/* ***Need to match the cards from the brands sites data***  */}
            {/* Will only have name of item and possibly price */}
        </>
    )

}