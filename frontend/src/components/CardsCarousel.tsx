
//female photos
import swimwear_model_v2_1 from "../assets/swimwear_v2_model1.jpg"
import swimwear_model_v2_2 from "../assets/swimwear_v2_model2.jpg"
import swimwear_model_v2_3 from "../assets/swimwear_v2_model3.jpg"
import swimwear_model_v2_4 from "../assets/swimwear_v2_model4.jpg"
import swimwear_model_v3_1 from "../assets/swimwear_v3_model1.jpg"
import swimwear_model_v3_2 from "../assets/swimwear_v3_model2.jpg"
import swimwear_model_v3_3 from "../assets/swimwear_v3_model3.jpg"
import swimwear_model_v4_1 from "../assets/swimwear_v4_model1.jpg"
import swimwear_model_v4_2 from "../assets/swimwear_v4_model2.jpg"
import summer_clothes_v2_model1 from "../assets/summer_clothes_v2_model1.jpg"
import summer_clothes_v2_model2 from "../assets/summer_clothes_v2_model2.jpg"
import summer_clothes_v2_model3 from "../assets/summer_clothes_v2_model3.jpg"

//male photos
import swimwear_v1_model_male1 from "../assets/swimwear_v1_model_male1.jpg"
import swimwear_v1_model_male2 from "../assets/swimwear_v1_model_male2.jpg"
import swimwear_v1_model_male3 from "../assets/swimwear_v1_model_male3.jpg"
import swimwear_v2_model_male1 from "../assets/swimwear_v2_model_male1.jpg"
import swimwear_v2_model_male2 from "../assets/swimwear_v2_model_male2.jpg"
import swimwear_v2_model_male3 from "../assets/swimwear_v2_model_male3.jpg"
import summer_clothes_v1_model_male1 from "../assets/summer_clothes_v1_model_male1.jpg"
import summer_clothes_v1_model_male2 from "../assets/summer_clothes_v1_model_male2.jpg"
import summer_clothes_v1_model_male3 from "../assets/summer_clothes_v1_model_male3.jpg"

import type { CardCarousel } from "../types/CardsCarouselType"

import arrow_left from "../assets/arrow-left.svg"
import arrow_right from "../assets/arrow-right.svg"

import { useRef, useState } from "react"

import "../styles/four-cards.css"
import { Link } from "react-router"

const setCardValues = (imgSrc: string, name: string, price: number) => {
    return {
        imgSrc: imgSrc,
        name: name,
        price: price
    }
}

export const CardsCarousel = () => {
    const[card1, setCard1] = useState<CardCarousel>(setCardValues(swimwear_model_v2_1, "Bikini", 30))
    const[card2, setCard2] = useState<CardCarousel>(setCardValues(swimwear_model_v2_2, "Bikini", 30))
    const[card3, setCard3] = useState<CardCarousel>(setCardValues(swimwear_model_v2_3, "Bikini", 30))
    const[card4, setCard4] = useState<CardCarousel>(setCardValues(swimwear_model_v2_4, "Bikini", 30))
    const[card5, setCard5] = useState<CardCarousel>(setCardValues(swimwear_model_v3_1, "Bikini", 30))
    const[card6, setCard6] = useState<CardCarousel>(setCardValues(swimwear_model_v3_2, "Bikini", 30))
    const[card7, setCard7] = useState<CardCarousel>(setCardValues(swimwear_model_v3_3, "Bikini", 30))
    const[card8, setCard8] = useState<CardCarousel>(setCardValues(swimwear_model_v4_1, "Bikini", 30))
    const[card9, setCard9] = useState<CardCarousel>(setCardValues(swimwear_model_v4_2, "Bikini", 30))
    const[card10, setCard10] = useState<CardCarousel>(setCardValues(summer_clothes_v2_model1, "Bikini", 30))
    const[card11, setCard11] = useState<CardCarousel>(setCardValues(summer_clothes_v2_model2, "Bikini", 30))

    //reference choice gender
    const womanChoice = useRef<HTMLHeadingElement>(null);
    const manChoice = useRef<HTMLHeadingElement>(null);

    //reference carousel functionalities
    const arrowLeftRef = useRef<HTMLDivElement>(null);
    const arrowRightRef = useRef<HTMLDivElement>(null);
    const carouselImagesRef = useRef<HTMLDivElement>(null);
    let carouselValue = 1450;
    let carouselIncrement = 0;


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
    const isClicked = (element: any) => {
        if(element.current){
            if(element.current.classList.contains("woman-choice-four-cards")){
                setCard1(setCardValues(swimwear_model_v2_1, "Bikini", 30));
                setCard2(setCardValues(swimwear_model_v2_2, "Bikini", 30))
                setCard3(setCardValues(swimwear_model_v2_3, "Bikini", 30))
                setCard4(setCardValues(swimwear_model_v2_4, "Bikini", 30))
                setCard5(setCardValues(swimwear_model_v3_1, "Bikini", 30))
                setCard6(setCardValues(swimwear_model_v3_2, "Bikini", 30))
                setCard7(setCardValues(swimwear_model_v3_3, "Bikini", 30))
                setCard8(setCardValues(swimwear_model_v4_1, "Bikini", 30))
                setCard9(setCardValues(swimwear_model_v4_2, "Bikini", 30))
                setCard10(setCardValues(summer_clothes_v2_model1, "Bikini", 30))
                setCard11(setCardValues(summer_clothes_v2_model2, "Bikini", 30))
            }else if(element.current.classList.contains("man-choice-four-cards")){
                setCard1(setCardValues(swimwear_v1_model_male1, "Bikini", 30));
                setCard2(setCardValues(swimwear_v1_model_male2, "Bikini", 30))
                setCard3(setCardValues(swimwear_v1_model_male3, "Bikini", 30))
                setCard4(setCardValues(swimwear_v2_model_male1, "Bikini", 30))
                setCard5(setCardValues(swimwear_v2_model_male2, "Bikini", 30))
                setCard6(setCardValues(swimwear_v2_model_male3, "Bikini", 30))
                setCard7(setCardValues(summer_clothes_v1_model_male1, "Bikini", 30))
                setCard8(setCardValues(summer_clothes_v1_model_male2, "Bikini", 30))
                setCard9(setCardValues(summer_clothes_v1_model_male3, "Bikini", 30))
            }
        }
    }

    const ClickArrowLeft = () => {
        if(arrowLeftRef.current){
            if(carouselIncrement > 0){
                if(carouselImagesRef.current){
                    carouselIncrement -= carouselValue
                    carouselImagesRef.current.style.right = `${carouselIncrement}px`
                }
            }
        }
    }

    const CliclArrowRight = () => {
        if(arrowRightRef.current){
            if(carouselIncrement < 2900){
                if(carouselImagesRef.current){
                    carouselIncrement += carouselValue;
                    carouselImagesRef.current.style.right = `${carouselIncrement}px`
                }
            }
        }
    }

    return (
        <>
        <div className="four-cards-component-container d-flex flex-column gap-4 align-items-center">
            <div className="male-female-option d-flex flex-row gap-4 w-75">
                <h1 ref={womanChoice} onClick={() => {ChoiceStyling(womanChoice, manChoice); isClicked(womanChoice)}} 
                className="woman-choice-four-cards fw-bold border-bottom border-dark border-4">Woman</h1>

                <h1 ref={manChoice} onClick={() => {ChoiceStyling(manChoice, womanChoice); isClicked(manChoice)}}
                className="man-choice-four-cards">Man</h1>
            </div>
            <div className="cards-arrows-container d-flex position-relative">
                <div onClick={ClickArrowLeft}
                className="left-arrow-button-container position-absolute z-2 h-50 d-flex justify-content-end align-items-end">
                    <div className="left-arrow-button bg-black" ref={arrowLeftRef}>
                        <img src={arrow_left} alt="arrow left icon" />
                    </div>
                </div>

                <div className="cards-container d-flex flex-row gap-2" ref={carouselImagesRef}>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card1.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card1.name}</p>
                            <p>{card1.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card2.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card2.name}</p>
                            <p>{card2.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card3.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card3.name}</p>
                            <p>{card3.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card4.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card4.name}</p>
                            <p>{card4.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card5.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card5.name}</p>
                            <p>{card5.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card6.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card6.name}</p>
                            <p>{card6.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card7.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card7.name}</p>
                            <p>{card7.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card8.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card8.name}</p>
                            <p>{card8.price}</p>
                        </div>
                    </div>
                     <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card9.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card9.name}</p>
                            <p>{card9.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card10.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card10.name}</p>
                            <p>{card10.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem"}}>
                        <img className="card-img-top rounded-0" src={card11.imgSrc} alt="Card image cap"/>
                        <div className="card-body">
                            <p className="card-text fw-bold fs-20">{card11.name}</p>
                            <p>{card11.price}</p>
                        </div>
                    </div>
                    <div className="card rounded-0 border-0" style={{"width": "22rem", "height": "30rem",}}>
                        <div className="card-body d-flex align-items-end justify-content-center">
                            <Link to={"/"} className="cards-carousel-more-btn btn btn-dark fw-bold" style={{"padding": "20px 30px", "fontSize": "20px"}}>
                                See More
                            </Link>
                        </div>
                    </div>
                </div>

                <div onClick={CliclArrowRight}
                className="right-arrow-button-container position-absolute z-1 w-100 h-50 d-flex justify-content-end align-items-end">
                    <div className="right-arrow-button bg-black" ref={arrowRightRef}>
                        <img src={arrow_right} alt="arrow right icon" />
                    </div>
                </div>
            </div>
            
        </div>
        </>
    )
}