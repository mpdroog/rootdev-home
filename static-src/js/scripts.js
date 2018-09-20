(function($){

	// Declare global variables
	var map,
		mapLatLng;

	$(document).ready( function(){
		// Add background image
		$('.left-wrap .bg').backstretch('images/bg.jpg');

		// Validate newsletter form
		$('<div class="spinner"><div class="square"></div><div class="square"></div><div class="square"></div><div class="square"></div></div>').hide().appendTo('.newsletter');
		$('<div class="success"></div>').hide().appendTo('.newsletter');
		$('#newsletter-form').validate({
			rules: {
				email: { required: true, email: true }
			},
			messages: {
				email: {
					required: 'Email address is required',
					email: 'Email address is not valid'
				}
			},
			errorElement: 'span',
			errorPlacement: function(error, element){
				error.appendTo(element.parent());
			},
			submitHandler: function(form){
				$(form).hide();
				$('.newsletter').find('.spinner').css({ opacity: 0 }).show().animate({ opacity: 1 });
				$.post($(form).attr('action'), $(form).serialize(), function(data){
					$('.newsletter').find('.spinner').animate({opacity: 0}, function(){
						$(this).hide();
						$('.newsletter').find('.success').show().html('<i class="icon ion-ios7-checkmark-outline"></i> Thank you for contacting!').animate({opacity: 1});
					});
				}).fail(function(res) {
					$('.newsletter').find('.spinner').animate({opacity: 0}, function(){
						$(this).hide();
						$('.newsletter').find('.success').show().html('<i class="icon ion-ios7-checkmark-outline"></i> Sorry, failed contacting.').animate({opacity: 1});
					});
				});
				return false;
			}
		});

		// Add tabs functionality to the right side content
		jgtContentTabs();

		// Load the Google Map object
		if ( $('#map-canvas').length > 0 ) {
			jgtGoogleMap();
		}
		
		// Set the minimum height for the right side and bind on resize or orientation change
		jgtMinHeight();
		$(window).bind('resize orientationchange', function(){
			jgtMinHeight();
		});

	});

	// Set the minimum height for the right side
	function jgtMinHeight(){
		var leftWrap = $('.left-wrap'),
			rightWrap = $('.right-wrap');

		if ( Modernizr.mq('only screen and (max-width: 1200px)') == true ) {
			rightWrap.css({ 'min-height': $(window).height() - leftWrap.height() });
		} else {
			rightWrap.removeAttr('style');
		}
	}

	// Add tabs functionality to the right side content
	function jgtContentTabs(){
		var tabsNav = $('#menu'),
			tabsWrap = $('#main');

		tabsWrap.find('.main-section:gt(0)').hide();
		tabsNav.find('li:first').addClass('active');
		tabsNav.find('a').click(function(e){
			tabsWrap.find('.main-section').hide();
			tabsWrap.find($($(this).attr('href'))).fadeIn(800);
			tabsNav.find('li').removeClass('active');
			$(this).parent().addClass('active');
			e.preventDefault();

                        if ( $('#map-canvas').length > 0 ) {
                            mymap._onResize();
                        }

			// Fix background
			$('.left-wrap .bg').backstretch('resize');

		});
	}

        var mymap = L.map('map-canvas', {attributionControl: false});
	// Create and initialize the Google Map object
	function jgtGoogleMap(){
            mymap.setView([52.662031, 4.817720], 15);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
              attribution: ''
            }).addTo(mymap);
            L.marker([52.662031, 4.817720]).addTo(mymap);
	}

})(jQuery);
