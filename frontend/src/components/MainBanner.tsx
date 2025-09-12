
//Needed to add the option of passing parameters to components
interface MainBannerSingle {
    imageSource: string;
}

interface MainBannerDouble {
    imageSource1: string;
    imageSource2: string;
}

export const MainBannerSingle: React.FC<MainBannerSingle> = ({imageSource}) => {

    return(
    <>
        {/* ***Loads a single image with the width of the screen/page*** */}
    </>
    )

}

export const MainBannerDouble: React.FC<MainBannerDouble> = ({imageSource1, imageSource2 }) => {

    return(
    <>
        {/* loads a two image carousel with only that changes the image every 10 seconds with the 
        width of the screen/page */}

        {/* Adding these here because otherwise typescript gets angry at me ( TypeScript => >:( ) */}
        <img src={imageSource1} alt="Banner 1" />
        <img src={imageSource2} alt="Banner 2" />
    </>
    )

}