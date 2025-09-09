import swimwear_model_v2_1 from "../assets/swimwear_v2_model1.jpg"
import swimwear_model_v2_2 from "../assets/swimwear_v2_model2.jpg"
import swimwear_model_v2_3 from "../assets/swimwear_v2_model3.jpg"
import swimwear_model_v2_4 from "../assets/swimwear_v2_model4.jpg"
import arrow_left from "../assets/arrow-left.svg"
import arrow_right from "../assets/arrow-right.svg"

import { useRef } from "react"

import "../styles/four-cards.css"

export const FourCards = () => {

    // This component is the predecessor of the cards carousel wheel and 
    // is kept here in case of any desire for reuse and resycle of the components structure. 

    const womanChoice = useRef(null);
    const manChoice = useRef(null);

    const ChoiceStyling = (choice1: any, choice2: any) => {
        if(choice1 != null && choice2 != null){
            choice1.current.classList.add("fw-bold");
            choice1.current.classList.add("border-bottom");
            choice1.current.classList.add("border-dark");
            choice1.current.classList.add("border-4");

            choice2.current.classList.remove("fw-bold");
            choice2.current.classList.remove("border-bottom");
            choice2.current.classList.remove("border-dark");
            choice2.current.classList.remove("border-4");
        }
    } 

    return (
        <>
        <div className="four-cards-component-container d-flex flex-column gap-4 align-items-center">
            <div className="male-female-option d-flex flex-row gap-4 w-75">
                <h1 ref={womanChoice} onClick={() => ChoiceStyling(womanChoice, manChoice)} 
                className="woman-choice-four-cards fw-bold border-bottom border-dark border-4">Woman</h1>

                <h1 ref={manChoice} onClick={() => ChoiceStyling(manChoice, womanChoice)}
                className="man-choice-four-cards">Man</h1>
            </div>
            <div className="cards-arrows-container d-flex position-relative">
                <div className="left-arrow-button-container position-absolute z-2 h-50 d-flex justify-content-end align-items-end">
                    <div className="left-arrow-button bg-black">
                        <img src={arrow_left} alt="arrow left icon" />
                    </div>
                </div>

                <div className="cards-container d-flex flex-row gap-2">
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={swimwear_model_v2_1} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">Bikini</p>
                            <p>Price</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={swimwear_model_v2_2} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">Bikini</p>
                            <p>Price</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={swimwear_model_v2_3} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">Bikini</p>
                            <p>Price</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={swimwear_model_v2_4} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">Bikini</p>
                            <p>Price</p>
                        </div>
                    </div>
                </div>

                <div className="right-arrow-button-container position-absolute z-1 w-100 h-50 d-flex justify-content-end align-items-end">
                    <div className="right-arrow-button bg-black">
                        <img src={arrow_right} alt="arrow right icon" />
                    </div>
                </div>
            </div>
            
        </div>
        </>
    )
}