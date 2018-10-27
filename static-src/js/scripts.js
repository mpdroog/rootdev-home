(function($){
	var map, mapLatLng;

	$(document).ready( function(){
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

		jgtContentTabs();
		if ( $('#map-canvas').length > 0 ) {
			jgtGoogleMap();
		}
	});

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
	function jgtGoogleMap(){
            mymap.setView([52.662031, 4.817720], 15);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
              attribution: ''
            }).addTo(mymap);
            L.marker([52.662031, 4.817720]).addTo(mymap);
	}

})(jQuery);
