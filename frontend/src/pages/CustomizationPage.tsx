import "../styles/customization-page.css"
import ShirtIcon from "../assets/shirt-icon.svg"
import PantsIcon from "../assets/pants-icon.svg"
import { NavBar } from "../partials/NavBar"

export const CustomizationPage = () => {
    return (
        <>
            <NavBar/>
            <div className="customization-page-container d-flex flex-column align-items-center gap-5">
                <div className="customization-text d-flex flex-column">
                    <div className="heading d-flex align-items-center">
                        <h1 className="fw-bold font-style-fjallaone fs-60">Customize</h1>
                    </div>
                    <div className="text">
                        <h2>
                            Take picking your style to another level <br />
                            with our customization options.
                        </h2>
                    </div>
                </div>
                <div className="customization-fields w-75 d-flex flex-column gap-3">
                    <div className="customize-upper-body card px-5 pt-3 pb-3 d-flex flex-column w-100 gap-3">
                        <div className="heading-upper-body d-flex gap-3">
                            <h2 className="fw-bold">Upper Body</h2>
                            <img src={ShirtIcon} alt="Shirt Icon"/>
                        </div>
                        <div className="upper-body-options d-flex justify-content-between">
                            <div className="neck-input d-flex flex-column">
                                <label htmlFor="neck">Neck Circumference</label>
                                <input type="text" name="neck" id="neck" />
                            </div>
                            <div className="chest-input d-flex flex-column">
                                <label htmlFor="chest">Chest Circumference</label>
                                <input type="text" name="chest" id="chest" />
                            </div>
                            <div className="torse-length-input d-flex flex-column">
                                <label htmlFor="torse-length">Tors Length</label>
                                <input type="text" name="torse-length" id="torse-length" />
                            </div>
                        </div>
                    </div>
                    <div className="customize-lower-body card px-5 pt-3 pb-3 d-flex flex-column w-100 gap-3">
                        <div className="heading-lower-body d-flex gap-3">
                            <h2 className="fw-bold">Lower Body</h2>
                            <img src={PantsIcon} alt="Pants Icon"/>
                        </div>
                        <div className="lower-body-options d-flex justify-content-between">
                            <div className="waist-input d-flex flex-column">
                                <label htmlFor="waist">Waist Circumference</label>
                                <input type="text" name="waist" id="waist" />
                            </div>
                            <div className="leg-len-input d-flex flex-column">
                                <label htmlFor="leg-len">Leg Length</label>
                                <input type="text" name="leg-len" id="leg-len" />
                            </div>
                            <div className="leg-wid-input d-flex flex-column">
                                <label htmlFor="leg-wid">Leg Width</label>
                                <input type="text" name="leg-wid" id="leg-wid" />
                            </div>
                        </div>
                    </div>
                    <div className="customize-clothing card px-5 pt-3 pb-3 d-flex flex-column w-100 gap-3">
                         <div className="heading-lower-body d-flex gap-3">
                            <h2 className="fw-bold">Additional</h2>
                            {/* <img src={} alt="Pants Icon"/> */}
                        </div>
                        <div className="additional-options d-flex justify-content-between">
                           <div className="color-input d-flex flex-column">
                                <label htmlFor="color">Clothing Color</label>
                                <input type="text" name="color" id="color" />
                            </div>
                            <div className="material-input d-flex flex-column">
                                <label htmlFor="material">Clothing material</label>
                                <input type="text" name="material" id="material" />
                            </div>
                            <div className="style-inpu d-flex flex-column">
                                <label htmlFor="style">Clothing Style</label>
                                <input type="text" name="style" id="style" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}