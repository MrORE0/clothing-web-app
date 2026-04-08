import type { ItemType } from "../types/CardsTypes"

interface FilteredContentInt {
    filteredContent: Array<ItemType>
}

export const FilteredContentPageDouble: React.FC<FilteredContentInt> = (filteredContent) => {

    return (
        <>

            {/* ***Renders filtered content here on a double column layout***  */}
            {/* ***Needs to have lazy loading mechanism***  */}

        </>
    )

}