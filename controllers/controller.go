package controllers

import (
	"net/http"

	movieModel "../models/movie"

	"github.com/gin-gonic/gin"
)

// CreateMovie parses request body and calls the movie model to create the movie
func CreateMovie(c *gin.Context) {
	var movie movieModel.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	_, err := movieModel.CreateMovie(&movie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// RetrieveMovie parses out the movie id parameter and calls the movie model to retrieve the movie
func RetrieveMovie(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required id parameter."})
		return
	}

	movie, err := movieModel.RetrieveMovie(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// UpdateMovie parses the request body and calls the movie model to update the movie
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required id parameter."})
		return
	}

	var movie movieModel.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	movie.MovieID = id
	_, err := movieModel.UpdateMovie(&movie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// DeleteMovie parses for the movie id parameter and calls the movie model to delete the movie
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required id parameter."})
		return
	}

	_, err := movieModel.DeleteMovie(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted."})
}
